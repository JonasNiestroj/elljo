context('elementRender', () => {
    beforeEach(() => {
        cy.exec('../main cypress/fixtures/element.jo cypress/fixtures/element_index.html element_output.html')
        cy.visit('./element_output.html')
    })

    it('renders hello world paragraph', () => {
        cy.get('p').should(($p) => {
            expect($p).to.have.length(1)
            expect($p).have.text('Hello world')
        })
    })
})
