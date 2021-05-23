describe('attributeRender', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/attribute.jo' })
        cy.document().invoke({ log: true }, 'write', html)
    })
    
    afterEach(() => {
        cy.visit("index.html")
    })


    it('renders hello world paragraph with title', () => {
        cy.get('p').should(($p) => {
          expect($p).to.have.length(1)
          expect($p).to.have.attr('title', 'hello world')
        })
    })
  
    it('updates title', () => {
        cy.get('button').click()
        cy.get('p').should(($p) => {
            expect($p).to.have.length(1)
            expect($p).to.have.attr('title', 'new world')
        })
    })
})
