import { Module } from '@nestjs/common';
import { HistoryTemplateService } from './history-template.service';
import { HistoryTemplateResolver } from './history-template.resolver';

@Module({
  providers: [HistoryTemplateService, HistoryTemplateResolver],
})
export class HistoryTemplateModule {}