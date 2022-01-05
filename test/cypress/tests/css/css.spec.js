import mount from '../../utils/mount';
import CSS from './css.jo';

describe('css', () => {
  beforeEach(() => {
    mount(CSS);
  });

  it('paragraphs should have css', () => {
    cy.get('p').should('have.css', 'color', 'rgb(255, 0, 0)');
  });
});
