describe('conditionalRender', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/conditional.jo' })
        cy.document().invoke({ log: true }, 'write', html)
    })

    afterEach(() => {
        cy.visit("index.html")
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

    it('shows on double click', () => {
        cy.get('button').click()
        cy.get('button').click()
        cy.get('p').should('have.length', 1)
    })
})
