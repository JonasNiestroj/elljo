import mount from '../../utils/mount';
import Attribute from './attribute.jo';

describe('attributeRender', () => {
  beforeEach(() => {
    mount(Attribute);
  });

  it('renders hello world paragraph with title', () => {
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).to.have.attr('title', 'hello world');
    });
  });

  it('updates title', () => {
    cy.get('button').click();
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).to.have.attr('title', 'new world');
    });
  });
});
