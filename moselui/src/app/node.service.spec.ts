import {
  beforeEachProviders,
  it,
  describe,
  expect,
  inject
} from '@angular/core/testing';
import { NodeService } from './node.service';

describe('NodeService Service', () => {
  beforeEachProviders(() => [NodeService]);

  it('should ...',
      inject([NodeService], (service: NodeService) => {
    expect(service).toBeTruthy();
  }));
});
