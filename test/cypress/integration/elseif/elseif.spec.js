describe('elseIfRender', () => {
    beforeEach(async () => {
        const html = await cy.task("spawnEllJo", { filePath: 'cypress/fixtures/elseif.jo' })
        cy.document().invoke({ log: true }, 'write', html)
    })

    afterEach(() => {
        cy.visit("index.html")
    })

    it('renders bye world paragraph', () => {
        cy.get('p').should(($p) => {
            expect($p).to.have.length(1)
            expect($p).have.text('Bye world')
        })
    })

    it('shows second world on click', () =>{
        cy.get('button').click()
        cy.get('p').should(($p) => {
          expect($p).to.have.length(1)
          expect($p).have.text('Second world')
        })
    })

  it('shows hello world on double click', () => {
      cy.get('button').click()
      cy.get('button').click()
      cy.get('p').should(($p) => {
          expect($p).to.have.length(1)
          expect($p).have.text('Hello world')
      })
    })
})
