import {NodeService} from "./node.service";
import {Component} from "@angular/core";
import {CHART_DIRECTIVES} from "angular2-highcharts/index";

@Component({
  moduleId: module.id,
  selector: 'moselui-app',
  templateUrl: 'moselui.component.html',
  styleUrls: ['moselui.component.css'],
  directives: [CHART_DIRECTIVES],
  providers: [NodeService],
})
export class MoseluiAppComponent {
  constructor(private nodeService: NodeService) {
    this.options = {
      title: {text: 'Test'},
      series: [{
        data: [29.9, 71.5, 106.4, 129.2],
      }]
    };
  }

  options:Object;
}
