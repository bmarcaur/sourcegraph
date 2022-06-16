import React, { useCallback } from 'react'

import InfoCircleOutlineIcon from 'mdi-react/InfoCircleOutlineIcon'

import { SettingsOrgSubject, SettingsUserSubject } from '@sourcegraph/shared/src/settings/settings'
import { Icon, Select } from '@sourcegraph/wildcard'

const getNamespaceDisplayName = (namespace: SettingsUserSubject | SettingsOrgSubject): string => {
    switch (namespace.__typename) {
        case 'User':
            return namespace.displayName ?? namespace.username
        case 'Org':
            return namespace.displayName ?? namespace.name
    }
}

const NAMESPACE_SELECTOR_ID = 'batch-spec-execution-namespace-selector'

interface NamespaceSelectorProps {
    namespaces: (SettingsUserSubject | SettingsOrgSubject)[]
    selectedNamespace: string
    disabled?: boolean
    onSelect: (namespace: SettingsUserSubject | SettingsOrgSubject) => void
}

export const NamespaceSelector: React.FunctionComponent<React.PropsWithChildren<NamespaceSelectorProps>> = ({
    namespaces,
    disabled,
    selectedNamespace,
    onSelect,
}) => {
    const onSelectNamespace = useCallback<React.ChangeEventHandler<HTMLSelectElement>>(
        event => {
            const selectedNamespace = namespaces.find(
                (namespace): namespace is SettingsUserSubject | SettingsOrgSubject =>
                    namespace.id === event.target.value
            )
            onSelect(selectedNamespace || namespaces[0])
        },
        [onSelect, namespaces]
    )

    return (
        <Select
            label={
                <>
                    <strong className="text-nowrap mb-2">Namespace</strong>
                    <Icon
                        aria-label="Coming soon"
                        data-tooltip="Coming soon"
                        as={InfoCircleOutlineIcon}
                        className="ml-1"
                    />
                </>
            }
            isCustomStyle={true}
            id={NAMESPACE_SELECTOR_ID}
            value={selectedNamespace}
            onChange={onSelectNamespace}
            disabled={disabled}
        >
            {namespaces.map(namespace => (
                <option key={namespace.id} value={namespace.id}>
                    {getNamespaceDisplayName(namespace)}
                </option>
            ))}
        </Select>
    )
}
