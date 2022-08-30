import { TextDocumentPositionParameters } from '@sourcegraph/client-api'
import { MaybeLoadingResult } from '@sourcegraph/codeintellify'
import * as comlink from 'comlink'
import { from, Subscription } from 'rxjs'
import { first } from 'rxjs/operators'
import { Unsubscribable } from 'sourcegraph'
import { newCodeIntelAPI } from '../../codeintel/api'

import { PlatformContext, ClosableEndpointPair } from '../../platform/context'
import { isSettingsValid } from '../../settings/settings'
import { FlatExtensionHostAPI, MainThreadAPI } from '../contract'
import { ExtensionHostAPIFactory } from '../extension/api/api'
import { proxySubscribable } from '../extension/api/common'
import { InitData } from '../extension/extensionHost'
import { registerComlinkTransferHandlers } from '../util'

import { ClientAPI } from './api/api'
import { ExposedToClient, initMainThreadAPI } from './mainthread-api'

export interface ExtensionHostClientConnection {
    /**
     * Closes the connection to and terminates the extension host.
     */
    unsubscribe(): void
}

/**
 * An activated extension.
 */
export interface ActivatedExtension {
    /**
     * The extension's extension ID (which uniquely identifies it among all activated extensions).
     */
    id: string

    /**
     * Deactivate the extension (by calling its "deactivate" function, if any).
     */
    deactivate(): void | Promise<void>
}

/**
 * @param endpoints The Worker object to communicate with
 */
export async function createExtensionHostClientConnection(
    endpointsPromise: Promise<ClosableEndpointPair>,
    initData: Omit<InitData, 'initialSettings'>,
    platformContext: Pick<
        PlatformContext,
        | 'settings'
        | 'updateSettings'
        | 'getGraphQLClient'
        | 'requestGraphQL'
        | 'telemetryService'
        | 'sideloadedExtensionURL'
        | 'getScriptURLForExtension'
        | 'clientApplication'
    >
): Promise<{
    subscription: Unsubscribable
    api: comlink.Remote<FlatExtensionHostAPI>
    mainThreadAPI: MainThreadAPI
    exposedToClient: ExposedToClient
}> {
    const subscription = new Subscription()

    // MAIN THREAD

    registerComlinkTransferHandlers()

    const { endpoints, subscription: endpointsSubscription } = await endpointsPromise
    subscription.add(endpointsSubscription)

    /** Proxy to the exposed extension host API */
    const initializeExtensionHost = comlink.wrap<ExtensionHostAPIFactory>(endpoints.proxy)

    const initialSettings = await from(platformContext.settings).pipe(first()).toPromise()
    const proxy = await initializeExtensionHost({
        ...initData,
        // TODO what to do in error case?
        initialSettings: isSettingsValid(initialSettings) ? initialSettings : { final: {}, subjects: [] },
    })

    const { api: newAPI, exposedToClient, subscription: apiSubscriptions } = initMainThreadAPI(proxy, platformContext)

    subscription.add(apiSubscriptions)

    const clientAPI: ClientAPI = {
        ping: () => 'pong',
        ...newAPI,
    }

    comlink.expose(clientAPI, endpoints.expose)
    proxy.mainThreadAPIInitialized().catch(() => {
        console.error('Error notifying extension host of main thread API init.')
    })

    // TODO(tj): return MainThreadAPI and add to Controller interface
    // to allow app to interact with APIs whose state lives in the main thread
    return { subscription, api: injectNewCodeintel(proxy), mainThreadAPI: newAPI, exposedToClient }
}

function injectNewCodeintel(old: comlink.Remote<FlatExtensionHostAPI>): comlink.Remote<FlatExtensionHostAPI> {
    const codeintel = newCodeIntelAPI({} as any)
    function thenMaybeLoadingResult<T>(promise: Promise<T>): Promise<MaybeLoadingResult<T>> {
        return promise.then<MaybeLoadingResult<T>>(result => {
            const maybeLoadingResult: MaybeLoadingResult<T> = { isLoading: false, result }
            return maybeLoadingResult
        })
    }

    const getHover = (textParameters: TextDocumentPositionParameters) => {
        console.log({ textParameters })
        return proxySubscribable(from(thenMaybeLoadingResult(codeintel.getHover(textParameters))))
    }
    return {
        ...old,
        getHover: getHover as any,
    }
}
