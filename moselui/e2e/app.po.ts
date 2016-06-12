export class MoseluiPage {
  navigateTo() {
    return browser.get('/');
  }

  getParagraphText() {
    return element(by.css('moselui-app h1')).getText();
  }
}
