import { EllJo } from './src/elljo';
import { afterRender, beforeDestroy, setComponent } from './src/lifecycleHooks';
import EllJoComponent from './src/elljoComponent';
import { createFragment } from './src/utils/fragment';

export default EllJo;
export {
  afterRender,
  beforeDestroy,
  setComponent,
  EllJoComponent,
  createFragment,
};
