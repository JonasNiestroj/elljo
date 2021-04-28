context('conditionalRender', () => {
    beforeEach(() => {
        cy.exec('../main cypress/fixtures/conditional.jo cypress/fixtures/conditional_index.html conditional_output.html')
        cy.visit('./conditional_output.html')
    })

    it('renders hello world paragraph', () => {
        cy.get('p').should(($p) => {
            expect($p).to.have.length(1)
            expect($p).have.text('Hello world')
        })
    })

    it('hides paragraph on click', () =>{
        cy.get('button').click()
        cy.get('p').should('have.length', 0)
    })

    it('shows on twice click', () => {
        cy.get('button').click()
        cy.get('button').click()
        cy.get('p').should('have.length', 1)
    })
})
