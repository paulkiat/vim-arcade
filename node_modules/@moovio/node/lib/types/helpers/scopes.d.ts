/**
 * Available scopes to request on OAuth tokens.
 */
export type SCOPES = any;
export namespace SCOPES {
    const ACCOUNTS_CREATE: string;
    const PROFILE_READ: string;
    const PROFILE_WRITE: string;
    const CONNECTIONS_READ: string;
    const CONNECTIONS_WRITE: string;
    const PING: string;
}
/**
 * For internal use only. Do not generate OAuth tokens for Moov.js and Moov
 * Drops that contain more permissions than are necessary.
 * @private
 */
export const ALL_SCOPES: string[];
//# sourceMappingURL=scopes.d.ts.map