const backendRootUrl = process.env.REACT_APP_BACKEND_ROOT_URL || 'localhost:6029';
const frontendRootUrl = process.env.REACT_APP_FRONTEND_ROOT_URL || 'localhost:3000';
const protocol = process.env.REACT_APP_PROTOCOL || 'http';

export const backendUrl = `${protocol}://${backendRootUrl}`;
export const frontendUrl = `${protocol}://${frontendRootUrl}`;
export const articles = ["der", "die", "das"];
export const defaultPageSize = 8