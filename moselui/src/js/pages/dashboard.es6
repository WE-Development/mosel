import {Controller} from "./controller.es6";

export class Dashboard extends Controller {

    constructor(container) {
        super(container, "pages/dashboard.html");
    }

    init() {
        console.log("Dashboard init");
        super.getChild("#chart")
            .highcharts({
                chart: {
                    type: 'bar'
                },
                title: {
                    text: 'Fruit Consumption'
                },
                xAxis: {
                    categories: ['Apples', 'Bananas', 'Oranges']
                },
                yAxis: {
                    title: {
                        text: 'Fruit eaten'
                    }
                },
                series: [{
                    name: 'Jane',
                    data: [1, 0, 4]
                }, {
                    name: 'John',
                    data: [5, 7, 3]
                }]
            });
    }

    destroy() {
        console.log("Dashboard destroy");
    }

}