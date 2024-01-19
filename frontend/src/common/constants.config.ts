export const API_BASE = "/api/v1";
export const BASE_URL = getBaseUrl();

export function getBaseUrl() {
  return "127.0.0.8080";
  // const baseUrl = window.mradis_config.host ?? process.env.REMARK_URL;

  // if (!baseUrl) {
  //   throw new Error(`Remark42: remark_config.host wasn't configured.`);
  // }

  // // Validate host
  // try {
  //   const { protocol } = new URL(baseUrl);
  //   // Show error if protocol of iframe doesn't match protocol of current page
  //   if (protocol !== window.location.protocol) {
  //     console.error('MRadis: Protocol mismatch.');
  //   }
  //   // Check if host has valid protocol and prevent XSS vurnuality
  //   if (!protocol.startsWith('http')) {
  //     console.error('MRadis: Wrong protocol in host URL.');
  //     throw new Error();
  //   }
  // } catch (e) {
  //   throw new Error('MRadis: Invalid host URL.');
  // }

  // return baseUrl;
}
