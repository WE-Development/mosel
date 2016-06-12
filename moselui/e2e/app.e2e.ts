import { MoseluiPage } from './app.po';

describe('moselui App', function() {
  let page: MoseluiPage;

  beforeEach(() => {
    page = new MoseluiPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('moselui works!');
  });
});
