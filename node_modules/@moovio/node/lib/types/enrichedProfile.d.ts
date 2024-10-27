/**
 * @typedef EnrichedProfile
 * @property {EnrichedIndividualProfile} individual - Describes a person
 * @property {EnrichedBusinessProfile} business - Describes a company
 *
 * @tag Enrichment
 */
/**
 * @typedef EnrichedBusinessProfile
 * @property {string} legalBusinessName
 * @property {EnrichedProfileAddress} address
 * @property {string} email
 * @property {EnrichedProfilePhone} phone
 * @property {EnrichedProfileIndustry} industryCodes - Describes industry specific identifiers
 * @property {string} website
 *
 * @tag Enrichment
 */
/**
 * @typedef EnrichedIndividualProfile
 * @property {EnrichedProfileName} name
 * @property {string} email
 * @property {EnrichedProfileAddress} address
 *
 * @tag Enrichment
 */
/**
 * @typedef EnrichedProfileAddress
 * @property {string} addressLine1
 * @property {string} addressLine2
 * @property {string} city
 * @property {string} stateOrProvince
 * @property {string} postalCode
 * @property {string} country
 *
 * @tag Enrichment
 */
/**
 * @typedef EnrichedProfileIndustry - Describes industry specific identifiers
 * @property {string} naics
 * @property {string} sic
 *
 * @tag Enrichment
 */
/**
 * @typedef EnrichedProfileName
 * @property {string} firstName
 * @property {string} middleName
 * @property {string} lastName
 * @property {string} suffix
 *
 * @tag Enrichment
 */
/**
 * @typedef EnrichedProfilePhone
 * @property {string} number
 * @property {string} countryCode
 *
 * @tag Enrichment
 */
/**
 * The Enriched Profile API.
 * @tag Enrichment
 */
export class EnrichedProfiles {
    constructor(moov: any);
    /**
     * @type {Moov}
     * @private
     */
    private moov;
    /**
     * Gets enriched profile data.
     *
     * @param {string} email - Email address associated with the profile.
     * @returns {Promise<EnrichedProfile>}
     * @tag Enrichment
     *
     * @example
     * const moov = new Moov(...);
     * try {
     *   const enrichedProfile = moov.enrichedProfiles.get("employee@business.com");
     * } catch (err) {
     *   // ..
     * }
     */
    get(email: string): Promise<EnrichedProfile>;
}
export type EnrichedProfile = {
    /**
     * - Describes a person
     */
    individual: EnrichedIndividualProfile;
    /**
     * - Describes a company
     */
    business: EnrichedBusinessProfile;
};
export type EnrichedBusinessProfile = {
    legalBusinessName: string;
    address: EnrichedProfileAddress;
    email: string;
    phone: EnrichedProfilePhone;
    /**
     * - Describes industry specific identifiers
     */
    industryCodes: EnrichedProfileIndustry;
    website: string;
};
export type EnrichedIndividualProfile = {
    name: EnrichedProfileName;
    email: string;
    address: EnrichedProfileAddress;
};
export type EnrichedProfileAddress = {
    addressLine1: string;
    addressLine2: string;
    city: string;
    stateOrProvince: string;
    postalCode: string;
    country: string;
};
/**
 * - Describes industry specific identifiers
 */
export type EnrichedProfileIndustry = {
    naics: string;
    sic: string;
};
export type EnrichedProfileName = {
    firstName: string;
    middleName: string;
    lastName: string;
    suffix: string;
};
export type EnrichedProfilePhone = {
    number: string;
    countryCode: string;
};
//# sourceMappingURL=enrichedProfile.d.ts.map