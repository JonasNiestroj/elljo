context('conditionalRender', () => {
    beforeEach(() => {
        cy.exec('../main cypress/fixtures/loop.jo cypress/fixtures/loop_index.html loop_output.html')
        cy.visit('./loop_output.html')
    })

    it('renders 3 children', () => {
        cy.get('p').should('have.length', 3)
    })

    it('adds one children on new array', () =>{
        cy.get('button').click()
        cy.get('p').should('have.length', 4)
    })
})
