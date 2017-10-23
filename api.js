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

class AliasAPI {
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
