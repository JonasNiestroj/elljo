describe('elementVariableRender', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/element_variable.jo' })
        cy.document().invoke({ log: true }, 'write', html)
    })
    
    afterEach(() => {
        cy.visit("index.html")
    })


    it('renders empty paragraph', () => {
        cy.get('p').should(($p) => {
            expect($p).to.have.length(1)
            expect($p).have.text('Hello world')
        })
    })

    it('changes text in paragraph', () => {
        cy.get('button').click()
        cy.get('p').should('have.text', 'Hello world!')
    })
})
