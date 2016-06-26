import {environment} from "./";
import { Injectable } from '@angular/core';
import {Http, Headers} from '@angular/http';

@Injectable()
export class NodeService {

  nodes: string[];

  constructor(private http: Http) {

    this.http.get(environment.baseUrl + '/info')
      //.map(res => res.text())
      .subscribe(
        data => console.log(data),
        err => console.error(err),
        () => console.log('Complete')
      );


  }

}
