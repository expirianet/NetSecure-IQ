const API_BASE_URL = process.env.VUE_APP_API_BASE_URL || '/api';

export default {
  // Auth
  LOGIN: `${API_BASE_URL}/login`,
  REGISTER: `${API_BASE_URL}/register`,
  
  // Organization
  ORGANIZATION: `${API_BASE_URL}/complete-organization`,
  ORGANIZATIONS: `${API_BASE_URL}/organizations`, // Pour les opérations CRUD complètes
  
  // Users
  USERS: `${API_BASE_URL}/users`,
  
  // MikroTik
  MIKROTIK: {
    RESOURCE: `${API_BASE_URL}/mikrotik/resource`,
    DUMP: `${API_BASE_URL}/mikrotik/dump`,
    PREREGISTER: `${API_BASE_URL}/mikrotik/preregister`,
    LIST: `${API_BASE_URL}/mikrotik/list`,
    TEST: `${API_BASE_URL}/mikrotik/test`,
    DISABLE: `${API_BASE_URL}/mikrotik/disable`,
    ENABLE: `${API_BASE_URL}/mikrotik/enable`,
    ASSOCIATE: `${API_BASE_URL}/mikrotik/associate`,
  },
  
  // Data
  ROUTERS: `${API_BASE_URL}/data/routers`,
  PING: `${API_BASE_URL}/ping`,
  
  // Protected route example
  PROTECTED: `${API_BASE_URL}/protected`,
};
