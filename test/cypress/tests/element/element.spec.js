import mount from '../../utils/mount';
import Element from './element.jo';

describe('elementRender', () => {
  beforeEach(() => {
    mount(Element);
  });

  it('renders hello world paragraph', () => {
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Hello world');
    });
  });
});
