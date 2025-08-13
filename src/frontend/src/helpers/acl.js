// src/helpers/acl.js

export function can(user, needed) {
    if (!user) return false;
    const perms = new Set((user.permissions || []).map(String));
    const list = Array.isArray(needed) ? needed : (needed ? [needed] : []);
    return list.every(p => perms.has(String(p)));
  }
  
  // Aussi un export default pour compatibilité
  export default can;
  