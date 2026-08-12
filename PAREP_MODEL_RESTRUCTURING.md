# PAREP Model Restructuring Summary

## Overview
The PAREP-CV data model has been restructured to align objectives, expected results, and indicators hierarchically. The table hierarchy was reorganized as follows:

## Key Changes

### 1. New Table: AMBITOS
- **Purpose**: Classify objectives by scope/area
- **Values**: 
  - `ambito_id: 1` = Pré-escolar (Objectives 1-4)
  - `ambito_id: 2` = 1º Ciclo (Objectives 5-7)

### 2. Table Renaming
| Old Name | New Name | Purpose |
|----------|----------|---------|
| `resultados` | `objetivos` | Strategic objectives of the program |
| `atividades` | `resultado_esperados` | Expected results that contribute to objectives |

### 3. Restructured Hierarchy
```
AMBITOS (Pré-escolar, 1º Ciclo)
  ├─ OBJETIVOS (7 total)
  │   ├─ RESULTADO_ESPERADOS (4-5 per objective, as bullets)
  │   └─ INDICADORES (with baseline and 2026 target metas)
  │       └─ METAS (targets for year 2026)
  ├─ ATIVIDADE_ANOS (annual execution)
  │   └─ ACOES (implementation actions)
  │       └─ CONTRATOS (contracts)
  │           ├─ CONTRATO_EXECUCOES (execution phases)
  │           └─ CONTRATO_ANEXOS (attachments)
  ├─ LOGISTICA (logistics management)
  └─ BENEFICIARIOS (beneficiary tracking)
```

### 4. Updated Foreign Keys
All tables now reference:
- `objetivo_id` instead of `resultado_id`
- `resultado_esperado_id` instead of `atividade_id`
- Preserved chain: objetivo → resultado_esperado → atividade_ano → acao → contrato

## Data Structure

### 7 Objectives organized by scope:

**Pré-escolar (ambito_id = 1)**
1. Assegurar acesso universal ao Pré-escolar
2. Melhorar desempenho das aprendizagens
3. Melhorar eficiência e eficácia dos recursos
4. Reforçar capacidade institucional e organizativa

**1º Ciclo (ambito_id = 2)**
5. Consolidar acesso equitativo e inclusivo
6. Reforçar êxito e qualidade das aprendizagens
7. Melhorar eficiência e eficácia dos recursos

### Expected Results Format
- Stored as text in `resultado_esperados` table
- 4-5 bulleted items per objective
- Examples:
  - "Acesso universal e inclusivo à Educação Pré-escolar"
  - "Redução das assimetrias regionais de acesso"

### Indicators & Metas
- **40 Total Indicators** (IND1-IND40)
- Each has:
  - `valor_baseline`: Baseline value (2021 or null if not available)
  - `unidades_id`: Unit of measurement (%, #, currency, days)
  - `meta_valor` for 2026 target in METAS table

## Files Modified

### 1. `parep_model_v2.md`
- Complete YAML schema restructuring
- Added AMBITOS table definition
- Renamed OBJETIVOS and RESULTADO_ESPERADOS
- Updated all foreign key references
- Updated INDICADORES to reference objetivo_id

### 2. `parep_data_objetivos.sql` (NEW)
- SQL INSERT statements for all 7 objectives
- All expected results data
- All 40 indicators with baseline values
- All 2026 target metas
- Organized by ambito (Pré-escolar, 1º Ciclo)

## Implementation Steps

1. **Update Database Schema**
   ```bash
   # Apply the new model structure
   # Run migrations to create AMBITOS table
   # Drop old constraints and recreate with new foreign keys
   ```

2. **Load Data**
   ```bash
   # Execute the SQL data file
   psql -U user -d database -f parep_data_objetivos.sql
   ```

3. **Migrate Existing Data** (if applicable)
   - Map old `resultados` records to new `objetivos`
   - Map old `atividades` records to new `resultado_esperados`
   - Update all foreign key references

## Benefits

✅ **Clearer Hierarchy**: Objectives → Expected Results → Indicators → Targets

✅ **Better Organization**: Scope-based classification (ambito) for filtering

✅ **Improved Semantics**: Table names now match the logical structure

✅ **Easier Reporting**: Expected results in bullets, indicators in metrics table

✅ **Consistent with Governance**: Aligns with PAREP-CV strategic framework
