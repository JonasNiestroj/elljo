import mount from '../../utils/mount';
import ElementVariable from './element_variable.jo';

describe('elementVariableRender', () => {
  beforeEach(() => {
    mount(ElementVariable);
  });

  it('renders empty paragraph', () => {
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Hello world');
    });
  });

  it('changes text in paragraph', () => {
    cy.get('button').click();
    cy.get('p').should('have.text', 'Hello world!');
  });
});
