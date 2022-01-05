import mount from '../../utils/mount';
import NoScript from './noscript.jo';

describe('noscript', () => {
  beforeEach(() => {
    mount(NoScript);
  });

  it('renders the paragraph', () => {
    cy.get('p').should('exist');
  });
});
