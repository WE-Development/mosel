import {environment} from "./";
import {Component} from "@angular/core";
import {CHART_DIRECTIVES} from "angular2-highcharts/index";

@Component({
  moduleId: module.id,
  selector: 'moselui-app',
  templateUrl: 'moselui.component.html',
  styleUrls: ['moselui.component.css'],
  directives: [CHART_DIRECTIVES],
})
export class MoseluiAppComponent {
  constructor() {
    this.options = {
      title: {text: environment.test},
      series: [{
        data: [29.9, 71.5, 106.4, 129.2],
      }]
    };
  }

  options:Object;
}
