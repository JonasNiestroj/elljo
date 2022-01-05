import mount from '../../utils/mount';
import ConditionalElse from './conditionalElse.jo';

describe('conditionalElseRender', () => {
  beforeEach(() => {
    mount(ConditionalElse);
  });

  it('renders hello world paragraph', () => {
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Hello world');
    });
  });

  it('shows paragraph on click', () => {
    cy.get('button').click();
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Bye world');
    });
  });

  it('shows if on double click', () => {
    cy.get('button').click();
    cy.get('button').click();
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Hello world');
    });
  });
});
