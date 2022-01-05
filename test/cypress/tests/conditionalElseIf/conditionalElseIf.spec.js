import mount from '../../utils/mount';
import ConditionalElseIf from './conditionalElseIf.jo';

describe('conditionalElseIfRender', () => {
  beforeEach(() => {
    mount(ConditionalElseIf);
  });

  it('renders bye world paragraph', () => {
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Bye world');
    });
  });

  it('shows second world on click', () => {
    cy.get('button').click();
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Second world');
    });
  });

  it('shows hello world on double click', () => {
    cy.get('button').click();
    cy.get('button').click();
    cy.get('p').should(($p) => {
      expect($p).to.have.length(1);
      expect($p).have.text('Hello world');
    });
  });
});
