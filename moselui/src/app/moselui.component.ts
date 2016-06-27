import {NodeService} from "./node.service";
import {Component} from "@angular/core";
import {CHART_DIRECTIVES} from "angular2-highcharts/index";
import {MD_CARD_DIRECTIVES} from '@angular2-material/card/card';
import {MD_GRID_LIST_DIRECTIVES} from "@angular2-material/grid-list/grid-list";

@Component({
  moduleId: module.id,
  selector: 'moselui-app',
  templateUrl: 'moselui.component.html',
  styleUrls: ['moselui.component.css'],
  directives: [
    CHART_DIRECTIVES,
    MD_CARD_DIRECTIVES,
    MD_GRID_LIST_DIRECTIVES
  ],
  providers: [NodeService],
})
export class MoseluiAppComponent {

  nodes = [];

  dogs: Object[] = [
    { name: 'Porter', human: 'Kara' },
    { name: 'Mal', human: 'Jeremy' },
    { name: 'Koby', human: 'Igor' },
    { name: 'Razzle', human: 'Ward' },
    { name: 'Molly', human: 'Rob' },
    { name: 'Husi', human: 'Matias' },
  ];

  constructor(private nodeService: NodeService) {
    this.nodeService.getNodes()
      .subscribe(
        res => {console.log(res);this.nodes = res.Nodes},
        err => console.log('Shit happens: ' + err))
  }

  options:Object;
}
