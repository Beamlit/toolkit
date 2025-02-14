import { Client as ModelContextProtocolClient } from "@modelcontextprotocol/sdk/client/index.js";
import { RunClient } from "../run.js";
import { Settings } from "../common/settings.js";
import { getSettings } from "../common/settings.js";
import { StructuredTool } from "@langchain/core/tools";
import { Client } from "@hey-api/client-fetch";
import { Function } from "../client/types.gen.js";
import { MCPClient, MCPToolkit } from "./mcp.js";
import { WebSocketClientTransport } from "./transport/websocket.js";
import { getAuthenticationHeaders } from "../authentication/authentication.js";
export type LocalFunction = {
    name: string;
    description: string;
    url: string;
};


/**
 * Toolkit for managing and interacting with remote toolkits and MCP services.
 */
export class LocalToolkit {
    private _client: Client;
    private modelContextProtocolClient: ModelContextProtocolClient;
    private url: string;
    private _functionName: string;
    private _function: Function | null = null;
    private _runClient: RunClient;
    private settings: Settings;

    /**
     * Creates an instance of RemoteToolkit.
     *
     * @param {Client} client - The HTTP client instance.
     * @param {string} functionName - The name of the remote function to manage.
     */
    constructor(client: Client, functionName: string, url: string) {
        this.settings = getSettings();
        this._client = client;
        this.modelContextProtocolClient = new ModelContextProtocolClient(
            {
                name: this.settings.name,
                version: "1.0.0",
            },
            {
                capabilities: {
                    tools: {}
                }
            }
        );

        this._functionName = functionName;
        this.url = url;
        this._runClient = new RunClient(client);
    }

    /**
     * Initializes the toolkit by retrieving the specified function and its associated tools.
     *
     * @returns {Promise<void>} Resolves when initialization is complete.
     * @throws Will throw an error if the function retrieval fails.
     */
    async initialize(name: string, description: string): Promise<void> {
        this._function = {
            metadata: {
                name: name,
            },
            spec: {
                integrationConnections: [],
            },
        };
        const headers = await getAuthenticationHeaders();
        const transport = new WebSocketClientTransport(new URL(this.url), {
            "x-beamlit-authorization": headers?.["X-Beamlit-Authorization"] || "",
            "x-beamlit-workspace": headers?.["X-Beamlit-Workspace"] || "",
        });
        await this.modelContextProtocolClient.connect(transport);
    }

    /**
     * Retrieves the list of structured tools from the remote function. If the function has integration connections,
     * it utilizes the MCPToolkit to manage them.
     *
     * @returns {Promise<StructuredTool[]>} An array of structured tools.
     * @throws Will throw an error if the toolkit has not been initialized.
     */
    async getTools(): Promise<StructuredTool[]> {
        if (!this._function) {
            throw new Error("Must initialize the toolkit first");
        }

        if (
            this._function.metadata &&
            this._function.spec?.integrationConnections
        ) {
            const mcpClient = new MCPClient(this.modelContextProtocolClient);
            const mcpToolkit = new MCPToolkit(mcpClient);
            await mcpToolkit.initialize();
            return mcpToolkit.getTools();
        }
        return [];
    }
}
