const fragmentCache = {};

export const createFragment = (html) => {
  let fragment = fragmentCache[html];
  if (!fragment) {
    let template = document.createElement('template');
    template.innerHTML = html;
    fragment = template.content;
    fragmentCache[html] = fragment;
  }
  return fragment.cloneNode(true);
};
