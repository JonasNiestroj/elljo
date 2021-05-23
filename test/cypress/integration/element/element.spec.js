describe('elementRender', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/element.jo' })
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
})
