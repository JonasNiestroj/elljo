import mount from '../../utils/mount';
import Svg from './svg.jo';

describe('svg', () => {
  beforeEach(() => {
    mount(Svg);
  });

  it('renders the svg', () => {
    cy.get('svg').should('exist');
  });
});
