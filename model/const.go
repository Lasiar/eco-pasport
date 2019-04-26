package model

const (
	sqlTest string = `SELECT 
org.org_name, 
org.Adress,  
t19.Allotted_wastewater_total, 
t19.Water_object, 
t11.Waste_generation_for_the_year, 
t8.Into_the_atmosphere, 
org.lat, 
org.lng 
from 
eco_2018.Table_0_5_Org org 
LEFT join ( 
select 
p1.Name, 
p2.Waste_generation_for_the_year 
from 
eco_2018.Table_1_11_part_1 p1 
inner join eco_2018.Table_1_11_part_2 p2 on 
p2.ID_p3 = p1.ID 
and p2.Hazard_class = 'всего' ) t11 on 
t11.Name = org.Org_name 
left join ( 
select 
pd.Allotted_wastewater_total, 
pd.Water_object, 
pd.Name 
from 
eco_2018.Table_1_9_Pollutant_discharges as pd ) t19 on 
org.Org_name = t19.name 
left join( 
select p1.Name, 
p2.Into_the_atmosphere 
from 
eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p1 p1 
inner join eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p2 p2 on 
p2.ID_p1 = p1.ID and p2.Name_of_pollutant = 'всего') t8 on 
t8.name = org.Org_name
where org.ID_Area = ?
order by org.Org_name
`

	sqlSpectial18 string = `select 
p1.Year,
p1.Economic_activ,
p1.Name,
p1.Emission_permit,
p2.Name_of_pollutant, 
p2.Thrown_without_cleaning_all, 
p2.Thrown_without_cleaning_organized, 
p2.Received_pollution_treatment, 
p2.Caught_and_rendered_harmless_all,
p2.Caught_and_rendered_harmless_utilized, 
p2.Into_the_atmosphere, 
p2.Sources_of_pollution_all, 
p2.Sources_of_pollution_organized, 
p2.MPE, 
p2.TAR,

p1.Source


from eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p1 p1
inner join eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p2 p2 on p1.ID = p2.ID_p1

where p1.ID_Area = ?`

	sqSpacial13 string = `SELECT
	[Year],
	Fee_total,
	Over_limit,
	From_stationary,
	Discharges,
	Waste_disposal,
	PNG,
	[Source]
FROM
	krasecology.eco_2018.Table_3_1_Fee_for_allowable_and_excess_emissions
where ID_Area = ?`

	sqlGetCenterArea string = `SELECT lat, lng from krasecology.eco_2018.Table_0_0_Regions where ID = ?`
	sqlGetInfoRegion string = `SELECT Admin_center , Creation_date, Population, Area, Gross_emissions, Withdrawn_water, Discharge_volume,Formed_waste  FROM eco_2018.Table_0_4_Regions_info WHERE Region_ID=?;`

	// TODO: переделать на уровне базы этот шлак
	sqlGetTableSpecial string = `select
	p1.[Year],
	p1.Economic_activity,
	p1.Name,
	p1.License,
	p1.Document_validity,
	p2.Standard,
	p2.Hazard_class,
	p2.Beginning_of_the_year,
	p2.Waste_generation_for_the_year,
	p2.Waste_receipt_all,
	p2.Waste_receipt_import,
	p2.Processed_waste,
	p2.Recycled_waste_all,
	p2.Recycled_waste_of_them_recycling,
	p2.Recycled_waste_of_them_processed,
	p2.Neutralized_all,
	p2.Neutralized_processed,
	p2.Waste_transfer_processing,
	p2.Waste_transfer_utilization,
	p2.Waste_transfer_neutralization,
	p2.Waste_transfer_storage,
	p2.Waste_transfer_burial,
	p2.Waste_disposal_storage,
	p2.Waste_disposal_burial,
	p2.End_of_the_year,
	p1.[Source]
FROM
	eco_2018.Table_1_11_part_1 p1
INNER JOIN eco_2018.Table_1_11_part_2 p2 on
	p2.ID_p3 = p1.ID
	and p2.ID_Area = ?`

	sqlGetTables string = "SELECT Table_ID, DB_Name, VisName FROM krasecology.eco_2018.Table_0_1_Tables"

	sqlGetRegions string = "SELECT id, num_region, name, cast(iif(is_town = 1,1,0) as BIT) from krasecology.eco_2018.Table_0_0_Regions"

	sqlGetHeaders string = `select
	'' as DB_Name,
	'' as VisName,
	header
from
	krasecology.eco_2018.Table_0_1_Tables
where
	header is not null
	and Table_ID = ?
union SELECT
	column_name,
	caption,
	null as header
from
	krasecology.eco_2018.Table_0_2_Columns
where Table_ID = ?
`

	sqlGetEmptyText string = "SELECT Empty_text FROM krasecology.eco_2018.Table_0_3_Empty_text where Table_ID = ? and Region_ID = ?"

	sqlGetSQL string = `
USE krasecology;

declare @SQL varchar(max) EXECUTE eco_2018.sp_get_table ?,@p1,
@p2,
@SQL output
EXECUTE (@sql)
`
)
