const { PrismaClient } = require('../schemes/prisma/generated/client/index')
const { randomUUID } = require('crypto')
const fs = require('fs')
const path = require('path')

/**
 * PrismaでSeedデータを追加するスクリプト
 *
 * 使い方:
 * 1. 事前にprisma generateを実行してクライアントを生成しておく
 *    $ pnpm prisma:generate
 * 2. このスクリプトを実行する
 *    $ pnpm db:seed
 */

const prisma = new PrismaClient()

/**
 * SQLファイルを読み込み、実行する
 * @param {string} filePath SQLファイルのパス
 */
async function executeSqlFile(filePath) {
  try {
    console.log(`Executing SQL file: ${filePath}`)
    const sql = fs.readFileSync(filePath, 'utf8')

    // SQLコマンドを分割して実行
    const queries = sql
      .replace(/(\r\n|\n|\r)/gm, ' ') // 改行を空白に置換
      .replace(/--.*$/gm, '') // コメントを削除
      .split(';') // セミコロンで分割
      .filter(query => query.trim() !== '') // 空のクエリを除外

    for (const query of queries) {
      await prisma.$executeRawUnsafe(query + ';')
    }

    console.log(`Successfully executed: ${filePath}`)
  } catch (error) {
    console.error(`Error executing SQL file ${filePath}:`, error)
    throw error
  }
}

/**
 * シードデータを追加する
 */
async function seed() {
  try {
    console.log('Starting database seeding...')

    // データベースのリセット（外部キー制約を一時的に無効化）
    const truncatePath = path.resolve(__dirname, '../schemes/seed/00_trancate.sql')
    await executeSqlFile(truncatePath)

    // シードデータの投入
    const seedPath = path.resolve(__dirname, '../schemes/seed/01_seed.sql')
    await executeSqlFile(seedPath)

    console.log('Database seeding completed successfully!')
  } catch (error) {
    console.error('Error seeding database:', error)
    process.exit(1)
  } finally {
    await prisma.$disconnect()
  }
}

// スクリプトの実行
seed()
