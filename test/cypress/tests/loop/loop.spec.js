import mount from '../../utils/mount';
import Loop from './loop.jo';

context('loop', () => {
  beforeEach(() => {
    mount(Loop);
  });

  it('renders 3 children', () => {
    cy.get('p').should('have.length', 3);
  });

  it('adds one children on new array', () => {
    cy.get('button').click();
    cy.get('p').should('have.length', 4);
  });
});
