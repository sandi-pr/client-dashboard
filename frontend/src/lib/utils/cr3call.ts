import axios, { type AxiosResponse } from 'axios';
interface Dict {
    [key: string]: any;
}

class Buffer {
    data: string;

    constructor(data: string) {
        this.data = data;
    }

    static fromJson(data: Dict): Buffer {
        return new Buffer(JSON.stringify(data));
    }
}

interface ConfigAPI {
    method: string;
    url_params: Options[];
    path: string;
}

interface ConfigServer {
    Port: number;
}

interface Config {
    Namespace?: string;
    Server?: ConfigServer;
    Templates?: ConfigTemplate[];
    APIs?: ConfigAPI[];
    Options?: Dict;
    source?: Buffer;
    file?: string;
    router?: any;
}

interface ConfigTemplateContext {
    Url?: string;
    Config?: Config;
    ConfigTemplate?: ConfigTemplate;
    URLParams?: Dict;
    Headers?: Dict;
    Request?: any;
    Writer?: any; 
}

interface ConfigTemplateCallback {
    Plugin?: string;
    Function?: string;
    plugin?: any;
    function?: (context: ConfigTemplateContext) => Dict; 
}

interface ConfigTemplate {
    Name?: string;
    Path?: string;
    FilePath?: string;
    Callback?: ConfigTemplateCallback;
    URLParams?: Options[];
    config?: Config;
}

interface Options {
    name: string;
    type: string;
}

export interface EndpointConfig {
    type?: string;
    url?: string;
    api?: ConfigAPI;
    template?: ConfigTemplate;
}

export interface CR3CallData {
    URLQuery: Dict;
    URLParams: Dict;
    Params: Dict;
    Headers: Dict;
}

function newDictFromJson(data: Buffer): Dict {
    const ret: Dict = data;
    return ret;
}


class CR3Call {
    private moduleManagerUrl: string;

    constructor(moduleManagerUrl: string) {
        this.moduleManagerUrl = moduleManagerUrl;
    }

    public async getConfig(moduleName: string, submoduleName: string): Promise<EndpointConfig | null> {
        try {
            const resp: AxiosResponse = await axios.get(`${this.moduleManagerUrl}/module-manager/get-config/${moduleName}/${submoduleName}`, {
                headers: { 'Content-Type': 'application/json' }
            });

            if (resp.status === 200) {
                const config: Dict = newDictFromJson(new Buffer(resp.data));
                const status: number | undefined = config.data.status;

                if (status !== undefined && status === 200) {
                    const result: Dict | undefined = config.data.result;

                    if (result !== undefined) {
                        const ret: EndpointConfig = {};
                        ret.api = result.api;
                        ret.template = result.template;
                        ret.type = result.type;
                        ret.url = result.url;

                        return ret;
                    }
                }
            }
        } catch (err) {
            console.error(err);
        }

        return null;
    }

    public async getPage(moduleName: string, submoduleName: string): Promise<string | null> {
        const config: EndpointConfig | null = await this.getConfig(moduleName, submoduleName);

        if (config !== null && config.type === 'page') {
            return config.url || null;
        }

        return null;
    }

    public async getAPI(moduleName: string, submoduleName: string): Promise<[string, ConfigAPI] | null> {
        const config: EndpointConfig | null = await this.getConfig(moduleName, submoduleName);

        if (config !== null && config.type === 'api' && config.api !== undefined) {
            return [config.url || '', config.api];
        }

        return null;
    }

    public async call(moduleName: string, submoduleName: string, data: CR3CallData): Promise<Dict | null> {
        const apiData: [string, ConfigAPI] | null = await this.getAPI(moduleName, submoduleName);

        if (apiData === null) {
            return null;
        }

        const [_tmp, api]: [string, ConfigAPI] = apiData;

        var url = this.moduleManagerUrl + '/' + moduleName + api.path;

        let headers: Dict = { 'Content-Type': 'application/json' };
        if (data.Headers !== null) {
            headers = { ...headers, ...data.Headers };
        }

        let configUrlParams: { [key: string]: Options } = {};

        if (api.url_params !== undefined && api.url_params !== null) {
            api.url_params.forEach((param) => {
                configUrlParams[param.name] = param;
            });
        }

        if (Object.keys(configUrlParams).length > 0 && data.URLParams !== null) {
            for (const [k, v] of Object.entries(data.URLParams)) {
                if (typeof v === 'string') {
                    url = url.replace(`:${k}`, v);
                } else if (typeof v === 'number') {
                    const configUrlParam: Options | undefined = configUrlParams[k];
                    if (configUrlParam !== undefined && (configUrlParam.type === 'number' || configUrlParam.type === 'integer' || configUrlParam.type === 'int')) {
                        url = url.replace(`:${k}`, v.toString());
                    } else if (configUrlParam !== undefined && (configUrlParam.type === 'float' || configUrlParam.type === 'double')) {
                        url = url.replace(`:${k}`, v.toFixed());
                    }
                }
            }
        }

        try {
            let resp: AxiosResponse;

            switch (api.method) {
                case 'GET':
                    resp = await axios.get(url, { params: data.URLQuery, data: data.Params, headers: headers });
                    break;
                case 'POST':
                    resp = await axios.post(url, data.Params, { params: data.URLQuery, headers: headers });
                    break;
                case 'PUT':
                    resp = await axios.put(url, data.Params, { params: data.URLQuery, headers: headers });
                    break;
                case 'DELETE':
                    resp = await axios.delete(url, { params: data.URLQuery, data: data.Params, headers: headers });
                    break;
                case 'PATCH':
                    resp = await axios.patch(url, data.Params, { params: data.URLQuery, headers: headers });
                    break;
                case 'HEAD':
                    resp = await axios.head(url, { params: data.URLQuery, data: data.Params, headers: headers });
                    break;
                case 'OPTIONS':
                    resp = await axios.options(url, { params: data.URLQuery, data: data.Params, headers: headers });
                    break;
                case 'TRACE':
                    resp = await axios.request({ method: 'TRACE', url: url, params: data.URLQuery, data: data.Params, headers: headers });
                    break;
                default:
                    throw new Error('Method not supported');
            }

            console.log(resp);

            if (resp.status === 200) {
                return newDictFromJson(resp.data);
            }
        } catch (err) {
            console.error(err);
        }

        return null;
    }

    public async getNamespaces(): Promise<string[] | null> {
        try {
            const resp: Dict | null = await this.call('module-manager', 'get-namespaces', { URLQuery: {}, URLParams: {}, Params: {}, Headers: {} });
            
            if (resp !== null) {
                const status: number | undefined = resp.status;
                if (status !== undefined && status === 200) {
                    const result: Dict | undefined = resp.result;
                    if (result !== undefined) {
                        const namespaces: string[] | undefined = result.namespaces;
                        if (namespaces !== undefined) {
                            return namespaces;
                        }
                    }
                }
            }
        } catch (err) {
            console.error(err);
        }

        return null;
    }

    public async getApis(moduleName: string): Promise<[string[], string[]] | null> {
        try {
            const resp: Dict | null = await this.call('module-manager', 'get-apis', {
                URLParams: { namespace: moduleName },
                URLQuery: {},
                Params: {},
                Headers: {}
            });

            if (resp !== null) {
                const status: number | undefined = resp.status;
                if (status !== undefined && status === 200) {
                    const result: Dict | undefined = resp.result;
                    if (result !== undefined) {
                        const apis: string[] | undefined = result.apis;
                        const pages: string[] | undefined = result.pages;
                        if (apis !== undefined && pages !== undefined) {
                            return [apis, pages];
                        }
                    }
                }
            }
        } catch (err) {
            console.error(err);
        }

        return null;
    }
}

export default CR3Call;