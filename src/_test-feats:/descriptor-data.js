// descriptorData.js

/**
 * @file descriptorData.js
 * @description Encapsulates data related to a connection or state within the application.
 */

/**
 * DescriptorData class to manage connection-related data.
 */
class DescriptorData {
  /**
   * Creates an instance of DescriptorData.
   * @param {Object} initialData - Initial data for the descriptor.
   */
  constructor(initialData = {}) {
    this.connectionId = initialData.connectionId || null;
    this.userId = initialData.userId || null;
    this.userName = initialData.userName || 'Guest';
    this.connectionStatus = initialData.connectionStatus || 'disconnected'; // e.g., connected, disconnected, connecting
    this.metadata = initialData.metadata || {}; // Additional data as needed
  }

  /**
   * Update the connection status.
   * @param {string} status 
   */
  updateConnectionStatus(status) {
    this.connectionStatus = status;
  }

  /**
   * Update user information.
   * @param {Object} userInfo 
   */
  updateUserInfo(userInfo = {}) {
    if (userInfo.userId !== undefined) this.userId = userInfo.userId;
    if (userInfo.userName !== undefined) this.userName = userInfo.userName;
  }

  /**
   * Update metadata.
   * @param {Object} newMetadata 
   */
  updateMetadata(newMetadata = {}) {
    this.metadata = { ...this.metadata, ...newMetadata };
  }

  /**
   * Get a summary of the descriptor data.
   * @returns {Object}
   */
  getSummary() {
    return {
      connectionId: this.connectionId,
      userId: this.userId,
      userName: this.userName,
      connectionStatus: this.connectionStatus,
      metadata: this.metadata
    };
  }
}

module.exports = DescriptorData;