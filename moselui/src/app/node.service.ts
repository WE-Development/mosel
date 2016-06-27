import {environment} from "./";
import { Injectable } from '@angular/core';
import {Http, Response} from '@angular/http';
import {Observable} from "rxjs/Observable";
import "rxjs/add/operator/map"

@Injectable()
export class NodeService {

  nodes: string[];

  constructor(private http: Http) {
  }

  getNodes(): Observable<any> {
    return this.http.get(environment.baseUrl + '/info')
      .map(this.extractData);
  }

  private extractData(res:Response) {
    let body = res.json();
    return body || { }
  }

}
