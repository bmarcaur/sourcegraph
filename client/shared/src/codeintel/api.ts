import { HoverMerged, TextDocumentPositionParameters } from '@sourcegraph/client-api'
import * as sourcegraph from './legacy-extensions/api'
import * as clientType from '@sourcegraph/extension-api-types'

export type QueryGraphQLFn<T> = () => Promise<T>

export interface CodeIntelAPI {
    hasReferenceProvidersForDocument(textParameters: TextDocumentPositionParameters): Promise<boolean>
    getDefinition(textParameters: TextDocumentPositionParameters): Promise<clientType.Location[]>
    getReferences(
        textParameters: TextDocumentPositionParameters,
        context: sourcegraph.ReferenceContext
    ): Promise<clientType.Location[]>
    getHover(textParameters: TextDocumentPositionParameters): Promise<HoverMerged>
    getDocumentHighlights(textParameters: TextDocumentPositionParameters): Promise<sourcegraph.DocumentHighlight[]>
}

export function newCodeIntelAPI(context: sourcegraph.CodeIntelContext): CodeIntelAPI {
    sourcegraph.updateCodeIntelContext(context)
    return new DefaultCodeIntelAPI()
}

class DefaultCodeIntelAPI implements CodeIntelAPI {
    hasReferenceProvidersForDocument(textParameters: TextDocumentPositionParameters): Promise<boolean> {
        return Promise.resolve(true)
    }
    getReferences(
        textParameters: TextDocumentPositionParameters,
        context: sourcegraph.ReferenceContext
    ): Promise<clientType.Location[]> {
        throw new Error('Method not implemented.')
    }
    getDefinition(textParameters: TextDocumentPositionParameters): Promise<clientType.Location[]> {
        return Promise.resolve([])
    }
    getHover(textParameters: TextDocumentPositionParameters): Promise<HoverMerged> {
        return Promise.resolve({ contents: [] })
    }
    getDocumentHighlights(textParameters: TextDocumentPositionParameters): Promise<sourcegraph.DocumentHighlight[]> {
        return Promise.resolve([])
    }
}
