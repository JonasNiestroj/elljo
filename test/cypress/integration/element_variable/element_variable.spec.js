context('elementVariableRender', () => {
    beforeEach(() => {
        cy.exec('../main cypress/fixtures/element_variable.jo cypress/fixtures/element_variable_index.html element_variable_output.html')
        cy.visit('./element_variable_output.html')
    })

    it('renders empty paragraph', () => {
        cy.get('p').should(($p) => {
            expect($p).to.have.length(1)
            expect($p).have.text('')
        })
    })

    it('changes text in paragraph', () => {
        cy.window().invoke('component.set', {text: 'Hello world'})
        cy.get('p').should('have.text', 'Hello world')
    })
})
