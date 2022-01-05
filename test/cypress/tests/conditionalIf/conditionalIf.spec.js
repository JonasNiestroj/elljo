import mount from '../../utils/mount';
import ConditionalIf from './conditionalIf.jo';

describe('conditionalIfRender', () => {
  beforeEach(() => {
    mount(ConditionalIf);
  });

  it('renders hello world paragraph', () => {
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Hello world');
    });
  });

  it('hides paragraph on click', () => {
    cy.get('button').click();
    cy.get('p').should('have.length', 0);
  });

  it('shows on double click', () => {
    cy.get('button').click();
    cy.get('button').click();
    cy.get('p').should('have.length', 1);
  });
});
