import { HoverMerged, TextDocumentIdentifier, TextDocumentPositionParameters } from '@sourcegraph/client-api'
import * as sourcegraph from './legacy-extensions/api'
import * as clientType from '@sourcegraph/extension-api-types'
import { languageSpecs } from './legacy-extensions/language-specs/languages'
import { LanguageSpec } from './legacy-extensions/language-specs/spec'
import { DocumentSelector, TextDocument } from 'sourcegraph'
import { match } from '../api/client/types/textDocument'
import { getModeFromPath } from '../languages'
import { parseRepoURI } from '../util/url'
import { createProviders, SourcegraphProviders } from './legacy-extensions/providers'
import { RedactingLogger } from './legacy-extensions/logging'
import { from, Observable } from 'rxjs'

export type QueryGraphQLFn<T> = () => Promise<T>

export interface CodeIntelAPI {
    hasReferenceProvidersForDocument(textParameters: TextDocumentPositionParameters): Promise<boolean>
    getDefinition(textParameters: TextDocumentPositionParameters): Promise<clientType.Location[]>
    getReferences(
        textParameters: TextDocumentPositionParameters,
        context: sourcegraph.ReferenceContext
    ): Promise<clientType.Location[]>
    getHover(textParameters: TextDocumentPositionParameters): Observable<HoverMerged>
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
    getHover(textParameters: TextDocumentPositionParameters): Observable<HoverMerged> {
        const x = findLanguageMatchingDocument(textParameters.textDocument)?.providers.hover.provideHover(
            {} as any,
            {} as any
        )
        if (!x) {
            return from(Promise.resolve({ contents: [] }))
        }
        x.forEach(y => {
            console.log({ y })
        })
        return from(Promise.resolve({ contents: [] }))
    }
    getDocumentHighlights(textParameters: TextDocumentPositionParameters): Promise<sourcegraph.DocumentHighlight[]> {
        return Promise.resolve([])
    }
}

function findLanguageMatchingDocument(textDocument: TextDocumentIdentifier): Language | undefined {
    const document: Pick<TextDocument, 'uri' | 'languageId'> = {
        uri: textDocument.uri,
        languageId: getModeFromPath(parseRepoURI(textDocument.uri).filePath || ''),
    }

    for (const language of languages) {
        if (match(language.selector, document)) {
            return language
        }
    }
    return undefined
}

interface Language {
    spec: LanguageSpec
    selector: DocumentSelector
    providers: SourcegraphProviders
}
const hasImplementationsField = true
const languages: Language[] = languageSpecs.map(spec => ({
    spec,
    selector: selectorForSpec(spec),
    providers: createProviders(spec, hasImplementationsField, new RedactingLogger(console)),
}))

function selectorForSpec(languageSpec: LanguageSpec): DocumentSelector {
    return [
        { language: languageSpec.languageID },
        ...(languageSpec.verbatimFilenames || []).flatMap(filename => [{ pattern: filename }]),
        ...languageSpec.fileExts.flatMap(extension => [{ pattern: `*.${extension}` }]),
    ]
}
