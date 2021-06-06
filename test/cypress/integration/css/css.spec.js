describe('css', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/css.jo' })
        cy.document().invoke({ log: true }, 'write', html)
    })
    
    afterEach(() => {
        cy.visit("index.html")
    })


    it('paragraphs should have css', () => {
      cy.get('p').should('have.css', 'color', 'rgb(255, 0, 0)')
    })
})
