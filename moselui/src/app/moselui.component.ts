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

  nodes = [];

  constructor(private nodeService: NodeService) {
    this.nodeService.getNodes()
      .subscribe(
        res => {console.log(res);this.nodes = res.Nodes},
        err => console.log('Shit happens: ' + err))
  }

  options:Object;
}
