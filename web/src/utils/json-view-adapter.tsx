import JsonView, { type JsonViewProps } from '@uiw/react-json-view';

export interface AdapterProps extends Omit<JsonViewProps<JsonViewProps<any>>, 'value'> {
    src: JsonViewProps<JsonViewProps<any>>['value'];        // 保留旧 prop 名
}

export function ReactJson({ src, ...rest }: AdapterProps) {
    return <JsonView value={src} {...rest} />;
}
