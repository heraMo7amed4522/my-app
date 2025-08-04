import { Module } from '@nestjs/common';
import { GraphQLModule } from '@nestjs/graphql';
import { ApolloDriver, ApolloDriverConfig } from '@nestjs/apollo';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { CardModule } from './card/card.module';
import { TransactionModule } from './transaction/transaction.module';
import { WalletModule } from './wallet/wallet.module';
import { PharaohModule } from './pharaoh/pharaoh.module';
import { HistoryTemplateModule } from './history-template/history-template.module';
import { QuizModule } from './quiz/quiz.module';

@Module({
  imports: [
    GraphQLModule.forRoot<ApolloDriverConfig>({
      driver: ApolloDriver,
      autoSchemaFile: true,
      playground: true,
      introspection: true,
    }),
    UserModule,
    AuthModule,
    CardModule,
    TransactionModule,
    WalletModule,
    PharaohModule,
    HistoryTemplateModule,
    QuizModule,
  ],
})
export class AppModule {}
