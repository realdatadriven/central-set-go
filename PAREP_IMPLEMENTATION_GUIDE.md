# PAREP Model v2 Restructuring - Complete Implementation Guide

## Executive Summary

The PAREP-CV database model has been successfully restructured to better align with the program's strategic framework. The key change repositions the data hierarchy from:

**Old Structure:**
```
Resultados → Atividades → Atividade_Anos → Ações
```

**New Structure:**
```
Âmbitos → Objetivos → Resultados Esperados → Atividade_Anos → Ações
```

This restructuring makes the model more intuitive and aligns expectations with the strategic planning framework used in education programs.

---

## What Changed

### 1. **New AMBITOS Table**
Introduces scope classification for all objectives and results.

```sql
CREATE TABLE ambitos (
  ambito_id SERIAL PRIMARY KEY,
  ambito VARCHAR(150) UNIQUE NOT NULL,
  ambito_desc TEXT,
  activo BOOLEAN DEFAULT true,
  -- timestamps
);

-- Values:
-- 1 = Pré-escolar (Educação Pré-escolar de qualidade, acessível e inclusiva)
-- 2 = 1º Ciclo (Ensino Básico de qualidade, equitativo, inclusivo e com êxito)
```

### 2. **Renamed Tables**

#### `resultados` → `objetivos`
- Now contains **7 strategic objectives** (previously 2 results)
- Each objective is associated with an `ambito_id`
- Focus: What the program aims to achieve

**Example Objectives:**
1. Assegurar que todas as crianças com idade de 4 e 5 anos frequentem o Pré-escolar
2. Melhorar o desempenho das aprendizagens das crianças no Pré-escolar
3. Melhorar a eficiência e eficácia do uso dos recursos ao Pré-escolar
4. Reforçar a capacidade institucional e organizativa
5. Consolidar o acesso equitativo e inclusivo no EBO
6. Reforçar o êxito e a qualidade das aprendizagens
7. Melhorar a eficiência e eficácia do uso dos recursos no EBO

#### `atividades` → `resultado_esperados`
- Now contains **expected results** (bullets that contribute to objectives)
- Each result ties to a specific objective
- Format: Descriptive text (4-5 per objective)

**Example Results for Objetivo 1:**
- Acesso universal e inclusivo à Educação Pré-escolar
- Redução das assimetrias regionais de acesso
- Aumento da frequência das crianças de 4 e 5 anos nos jardins
- Maior inclusão de crianças de famílias vulneráveis

### 3. **Updated Foreign Key References**

| Table | Old FK | New FK | Changed From | Changed To |
|-------|--------|--------|--------------|------------|
| indicadores | resultado_id | objetivo_id | Objectives 2 | Objectives 7 |
| atividade_anos | resultado_id | - | Removed | - |
| atividade_anos | atividade_id | resultado_esperado_id | Activities → Expected Results | - |
| acoes | resultado_id | objetivo_id | - | New |
| acoes | atividade_id | resultado_esperado_id | - | New |
| contratos | resultado_id | objetivo_id | - | New |
| contratos | atividade_id | resultado_esperado_id | - | New |
| logistica | resultado_id | objetivo_id | - | New |
| logistica | atividade_id | resultado_esperado_id | - | New |
| beneficiarios | resultado_id | objetivo_id | - | New |
| beneficiarios | atividade_id | resultado_esperado_id | - | New |

### 4. **New Data Structure**

#### 7 Objectives with Indicators

| Âmbito | Obj # | Objetivo | # Esperados | # Indicadores |
|--------|-------|----------|-------------|---------------|
| Pré-Escolar | 1 | Acesso Universal | 4 | 4 (IND1-4) |
| Pré-Escolar | 2 | Desempenho Aprendizagens | 4 | 5 (IND5-9) |
| Pré-Escolar | 3 | Eficiência Recursos | 3 | 2 (IND10-11) |
| Pré-Escolar | 4 | Capacidade Institucional | 5 | 9 (IND12-20) |
| 1º Ciclo | 5 | Acesso Equitativo | 4 | 6 (IND21-26) |
| 1º Ciclo | 6 | Êxito e Qualidade | 5 | 8 (IND27-34) |
| 1º Ciclo | 7 | Eficiência Recursos | 5 | 6 (IND35-40) |

#### 40 Indicators with Targets

Each indicator has:
- **Descritor**: Clear description
- **Baseline**: Value from 2021 (or null if unavailable)
- **Meta 2026**: Target value for 2026
- **Unidade**: Unit of measurement (%, #, etc.)

---

## Files Modified/Created

### 1. Modified: `assets/models/parep_model_v2.md`
**What Changed:**
- Added AMBITOS table definition
- Renamed RESULTADOS → OBJETIVOS table
- Renamed ATIVIDADES → RESULTADO_ESPERADOS table
- Updated all foreign key relationships
- Cleaned up old data section (moved to SQL file)
- Updated table order in menu/display

**Key Additions:**
```yaml
## AMBITOS Table
- ambito_id: 1, ambito: "Pré-escolar"
- ambito_id: 2, ambito: "1º Ciclo"

## OBJETIVOS Table
- New fk: ambito_id (foreign key to ambitos)

## RESULTADO_ESPERADOS Table
- New fk: objetivo_id (foreign key to objetivos)
```

### 2. New: `assets/models/parep_data_objetivos.sql`
**Purpose:** SQL INSERT statements for all objectives, results, indicators, and 2026 targets

**Contents:**
```
- 7 Objetivo INSERT statements (objectives)
- 28 Resultado_Esperado INSERT statements (expected results)
- 40 Indicador INSERT statements (indicators with baseline values)
- 40 Meta INSERT statements (2026 targets)
```

**Organized by:**
- Objectives 1-4: Pré-escolar
- Objectives 5-7: 1º Ciclo

### 3. New: `PAREP_MODEL_RESTRUCTURING.md`
Comprehensive guide documenting all changes, benefits, and implementation steps.

### 4. New: `PAREP_OBJECTIVES_REFERENCE.md`
Quick reference showing all 7 objectives, their expected results, and indicators in table format.

---

## Implementation Checklist

### Phase 1: Database Schema Update
- [ ] Back up current database
- [ ] Review and validate new YAML schema in `parep_model_v2.md`
- [ ] Create AMBITOS table (primary key: ambito_id)
- [ ] Rename RESULTADOS table to OBJETIVOS
- [ ] Add ambito_id foreign key to OBJETIVOS
- [ ] Rename ATIVIDADES table to RESULTADO_ESPERADOS
- [ ] Update ATIVIDADE_ANOS foreign keys:
  - Remove atividade_id foreign key
  - Add resultado_esperado_id foreign key
- [ ] Update ACOES table foreign keys (add objetivo_id, resultado_esperado_id)
- [ ] Update CONTRATOS table foreign keys (add objetivo_id, resultado_esperado_id)
- [ ] Update LOGISTICA table foreign keys (add objetivo_id, resultado_esperado_id)
- [ ] Update BENEFICIARIOS table foreign keys (add objetivo_id, resultado_esperado_id)
- [ ] Update INDICADORES table:
  - Drop resultado_id foreign key
  - Add objetivo_id foreign key
  - Validate all 40 indicators exist and updated

### Phase 2: Data Migration
- [ ] Migrate existing data from old RESULTADOS to new OBJETIVOS structure
- [ ] Migrate existing data from old ATIVIDADES to new RESULTADO_ESPERADOS
- [ ] Update all foreign key references in dependent tables
- [ ] Execute `parep_data_objetivos.sql` to load 2026 strategic targets
- [ ] Validate data integrity

### Phase 3: Verification
- [ ] Verify all 7 objectives loaded correctly
- [ ] Confirm all 28 expected results linked to objectives
- [ ] Validate all 40 indicators with baseline and 2026 metas
- [ ] Test reporting views/dashboards against new structure
- [ ] Update API endpoints to reflect new table names

### Phase 4: Documentation & Rollout
- [ ] Update API documentation
- [ ] Update UI/UX components to use new table names
- [ ] Train users on new structure
- [ ] Update backups and disaster recovery procedures
- [ ] Create migration rollback procedures

---

## Data Statistics

### Objectives Distribution
- **Pré-Escolar**: 4 objectives, 20 indicators
- **1º Ciclo**: 3 objectives, 20 indicators
- **Total**: 7 objectives, 40 indicators

### Indicators by Unit Type
- **Percentage (%)**: 17 indicators
- **Number (#)**: 23 indicators
- **Other**: 0 indicators

### Indicator Coverage
- **With Baseline**: 22 indicators (55%)
- **Without Baseline (null)**: 18 indicators (45%)

---

## Key Benefits

✅ **Strategic Alignment**: Structure aligns with official PAREP-CV strategic framework

✅ **Better Organization**: Scope-based classification enables filtering by Pré-Escolar or 1º Ciclo

✅ **Clearer Hierarchy**: Objectives → Results → Indicators → Targets (instead of flat structure)

✅ **Semantic Clarity**: Table names now reflect their business meaning

✅ **Flexible Framework**: Can easily extend to 3rd, 4th+ cycles in future

✅ **Reporting Ease**: Simplified queries for dashboard and report generation

✅ **Data Validation**: FK relationships ensure data consistency

---

## Backward Compatibility

⚠️ **Breaking Changes:**
- Old table names (`resultados`, `atividades`) are retired
- Applications using old table names must be updated
- API endpoints must be updated to reference new names

✅ **Mitigation:**
- Provide migration scripts for existing data
- Create database views as aliases (optional compatibility layer)
- Comprehensive documentation available

---

## Support & Questions

For questions or issues:
1. Review `PAREP_MODEL_RESTRUCTURING.md` for detailed changes
2. Check `PAREP_OBJECTIVES_REFERENCE.md` for data structure
3. Consult `parep_model_v2.md` for YAML schema details
4. Review `parep_data_objetivos.sql` for data format

---

## Appendix: SQL Migration Example

```sql
-- Example: If migrating existing data
BEGIN TRANSACTION;

-- 1. Create new AMBITOS
INSERT INTO ambitos (ambito_id, ambito, ambito_desc)
VALUES 
  (1, 'Pré-escolar', 'Educação Pré-escolar de qualidade'),
  (2, '1º Ciclo', 'Ensino Básico de qualidade');

-- 2. Migrate RESULTADOS → OBJETIVOS
ALTER TABLE resultados RENAME TO objetivos;
ALTER TABLE objetivos ADD COLUMN ambito_id INTEGER REFERENCES ambitos;

-- 3. Assign ambito_id based on business logic
UPDATE objetivos SET ambito_id = 1 WHERE objetivo_id IN (1,2,3,4);
UPDATE objetivos SET ambito_id = 2 WHERE objetivo_id IN (5,6,7);

-- 4. Add NOT NULL constraint
ALTER TABLE objetivos ALTER COLUMN ambito_id SET NOT NULL;

-- 5. Migrate ATIVIDADES → RESULTADO_ESPERADOS
ALTER TABLE atividades RENAME TO resultado_esperados;
ALTER TABLE resultado_esperados DROP COLUMN resultado_id;
ALTER TABLE resultado_esperados ADD COLUMN objetivo_id INTEGER NOT NULL REFERENCES objetivos;

-- 6. Update dependent tables
ALTER TABLE atividade_anos DROP CONSTRAINT atividade_anos_atividade_id_fkey;
ALTER TABLE atividade_anos ADD COLUMN resultado_esperado_id INTEGER REFERENCES resultado_esperados;
-- ... (similar for other tables)

COMMIT;
```

---

**Version:** 2.0  
**Date:** August 12, 2026  
**Status:** Complete
