export default (component) => {
  document.body.innerHTML = '';
  new component({ target: document.body });
};
