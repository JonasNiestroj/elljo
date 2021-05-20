context('loop', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/loop.jo' })
        cy.document().invoke({ log: true }, 'write', html)
    })
    
    afterEach(() => {
        cy.visit("index.html")
    })


    it('renders 3 children', () => {
        cy.get('p').should('have.length', 3)
    })

    it('adds one children on new array', () =>{
        cy.get('button').click()
        cy.get('p').should('have.length', 4)
    })
})
