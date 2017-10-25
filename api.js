/**
 * JS class to interact with aliasd
 * ================================
 *
 * Methods:
 * ========
 *
 * get(scope, id)      => returns the alias
 * getAll(scope)       => returns mapping from id to alias for a given scope
 * create(scope, id)   => create a new, random alias
 * set(scope, id, val) => set a alias
 * del(scope, id)      => delete a alias
 */
export class AliasAPI {
  constructor(endpoint){
    this.endpoint = endpoint;
  }

  get(scope, id) {
    return fetch(this.endpoint+`/alias/${scope}/${id}`, {
      method: 'GET',
      mode: 'cors',
    })
    .then((resp) => {
      if(resp.status !== 200){
        return resp.text().then((e)=>Promise.reject(e));
      }
      return resp.text();
    });
  }

  getAll(scope) {
    return fetch(this.endpoint+`/alias/${scope}`, {
      method: 'GET',
      mode: 'cors',
    })
    .then((resp) => {
      if(resp.status !== 200){
        return resp.text().then((e)=>Promise.reject(e));
      }
      return resp.json();
    });
  }

  create(scope, id) {
    return fetch(this.endpoint+`/alias/${scope}/${id}`, {
      method: 'POST',
      mode: 'cors',
    })
    .then((resp) => {
      if(resp.status !== 200){
        return resp.text().then((e)=>Promise.reject(e));
      }
      return resp.text();
    });
  }

  set(scope, id, value) {
    return fetch(this.endpoint+`/alias/${scope}/${id}`, {
      method: 'PUT',
      mode: 'cors',
      body: value,
    })
    .then((resp) => {
      if(resp.status !== 200){
        return resp.text().then((e)=>Promise.reject(e));
      }
      return resp.text();
    });
  }

  del(scope, id) {
    return fetch(this.endpoint+`/alias/${scope}/${id}`, {
      method: 'DELETE',
      mode: 'cors',
    })
    .then((resp) => {
      if(resp.status !== 200){
        return resp.text().then((e)=>Promise.reject(e));
      }
      return resp.text();
    });
  }
}


/**
 * JS class to interact with aliasd and cache the results
 * ================================
 *
 * Methods:
 * ========
 *
 * initScopes(scopes)  => initialize the given scopes (maps to getAll)
 * get(scope, id)      => returns the alias
 * create(scope, id)   => create a new, random alias
 * set(scope, id, val) => set a alias
 * del(scope, id)      => delete a alias
 */
export class AliasCache extends AliasAPI {
  constructor(url){
    super(url);
    this.scopes = {};
  }

  initScopes(scopes) {
    const promises = [];
    if (scopes !== undefined && scopes.length > 0) {
      for(let i=0; i<scopes.length; i++) {
        promises.push(this.getAll(scopes[i]));
      }
    }
    return Promise.all(promises).then((results) => {
      for(let i=0; i<scopes.length; i++) {
        const scope = scopes[i];
        this.scopes[scope] = results[i];
      }
      return this;
    });
  }

  get(scope, id) {
    return new Promise((resolve,reject)=>{
      const cache = this.scopes[scope] || {};
      if (cache[id] === undefined) {
        return super.get(scope, id)
        .then((alias)=>{
          resolve(alias);
        })
        .catch(reject);
      }else {
        resolve(cache[id]);
      }
    });
  }

  create(scope, id) {
    return super.create(scope, id)
    .then((alias) => {
      const cache = this.scopes[scope] || {};
      cache[id] = alias;
      this.scopes[scope] = cache;
      return alias;
    });
  }

  del(scope, id) {
    return super.del(scope, id)
    .then(() => {
      const cache = this.scopes[scope] || {};
      delete cache[id];
      this.scopes[scope] = cache;
      return true;
    });
  }

  set(scope, id, alias) {
    return super.set(scope, id, alias)
    .then(() => {
      const cache = this.scopes[scope] || {};
      cache[id] = alias;
      this.scopes[scope] = cache;
      return alias;
    });
  }

}
