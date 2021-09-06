SELECT lower(table_name), 'SELECT '||LISTAGG (output,'||')   WITHIN GROUP (ORDER BY table_name )  ||'from '||owner||'.' || table_name
	FROM 
	(
	select   owner,table_name, null_check_pre|| pre || mid || post || null_check_post  
			AS output
	from (
	select table_name,owner, 
		(select count(*) from dba_tab_columns c where c.owner=t.owner and c.table_name=t.table_name and c.column_name = 'FAD_HASH') as fh_chk, -- can be used to split by fad hash
		lead(column_name) over (partition by table_name order by column_id) x, -- if last column then do not concatenate ||
		lag(column_name) over (partition by table_name order by column_id) y, -- not used...
		column_name,
		data_type,        
		case nullable 
			when 'N' then ' '
			when 'Y' then ' coalesce('
			end as null_check_pre ,
	
			case when data_type in ('DATE', 'TIMESTAMP(4)', 'NUMBER') then
			'  trim(to_char( '
		 when data_type in ('CLOB') THEN
		'  to_char( '
		when data_type in ('BLOB') then
			' DBMS_OBFUSCATION_TOOLKIT.MD5(input=>'
		when data_type in ('VARCHAR2') and data_length > 1999 then  -- do we need to substr < 2000 for long chars?
			--' DBMS_OBFUSCATION_TOOLKIT.MD5(input=>UTL_RAW.cast_to_raw (COALESCE (SUBSTR ('
			--' substr( ' 
			--' upper(md5( coalesce( left('
			--' upper( left( '
			'substr('
		else ' '
		end as pre,
		
		column_name as mid,
		
		--select * from dba_tab_columns where owner = 'OPS$BRDB' and data_type = 'NUMBER' and data_scale > 0;
		case when data_type in ('DATE', 'TIMESTAMP(4)') then
			' , ''YYYYMMDDHH24MISS''))'
		when data_type = 'NUMBER' then
			case when column_name = 'JOURNAL_SEQ_NUMBER' then
				',''09999999999999''))'
			when table_name = 'BRDB_TC_MODES_MAPPING' and column_name = 'ALLOWED_MODES_ID' then
				',''09''))'
			when data_precision is null and data_scale is null then
				',''0' || lpad('9', 22, '9') || '''))'
			else
			',''0' || lpad('9', data_precision-data_scale, '9') || 
				case when data_scale > 0 then
				'.' || lpad('9', data_scale,'9')
				end
			|| '''))'
			end
			/*
			case when column_name = 'FAD_HASH' then
				', ''099''))'
			when column_name in ( 'BRANCH_ACCOUNTING_CODE', 'BRANCH_CODE') then
			', ''099999''))'
			when 
			else
			', ''0999999999999.99''))'
			end*/
		when data_type in ('BLOB','CLOB') then
			') '
		when data_type in ('VARCHAR2') and data_length > 1999 then
			--',1, 2000)' --, ''.''))) '
			--' ' --',2000),''.'')))'
			',1,2000) '
		else ' '
		end as post,
		case nullable 
			when 'N' then ' '
			when 'Y' then ' ,'' '')'
			end as null_check_post
	from   dba_tab_columns t
	where data_type not in ('BLOB','CLOB')
	and  (
	owner = 'BCMS' and table_name in (
		'IMPORT_CSV_BULK_FILE', 
	'IMPORT_CSV_BULK_FILE_AUDIT', 
	'OBC_ACTION_REQUEST', 
	'OBC_BRANCH_TYPE', 
	'OBC_DEVICE_TYPE', 
	'OBC_EXTERNAL_SYSTEM', 
	'OBC_FINANCIAL_YEAR', 
	'OBC_HIH_NODE_ID', 
	'OBC_INTEGRATOR', 
	'OBC_INTERFACE', 
	'OBC_INTERFACE_AUDIT', 
	'OBC_RESOURCE_META', 
	'OBC_RETAILER', 
	'OBC_STOCK_UNIT', 
	'OBC_STOCK_UNIT_TYPE', 
	'PATCH_RELEASE_MANAGEMENT'
	)
	)
	OR (
	owner = 'PAF_OWNER' and table_name in (
		'PAF_ADDRESS_POINT_A', 
		'PAF_ADDRESS_POINT_B'
	)
	)
	OR (
	owner = 'OPS$BRDBTR' and table_name in (
		'TEM_BRANCH_STOCK_UNITS', 
		'TEM_BRANCH_USERS', 
		'TEM_BRANCH_USER_ROLES', 
		'TEM_RX_BTS_DATA', 
		'TEM_STOCK_UNIT_ASSOCIATIONS', 
		'TEM_SU_OPENING_BALANCE', 
		'TEM_TXN_ACK_DETAILS'
	)
	)
	OR (
	owner = 'EMDB2' and table_name in (
		'EMDB2_ASSOCIATED_STOCK_UNIT', 
		'EMDB2_ASSO_STOCK_UNIT_AUDIT', 
		'EMDB2_BRANCH_NODES', 
		'EMDB2_BRANCH_NODES_AUDIT', 
		'EMDB2_POST_OFFICE', 
		'EMDB2_POST_OFFICE_AUDIT', 
		'EMDB2_RESOURCE_POOL', 
		'EMDB2_STOCK_UNIT', 
		'EMDB2_STOCK_UNIT_AUDIT', 
		'EMDB2_USERS', 
		'EMDB2_USERS_AUDIT', 
		'PATCH_RELEASE_MANAGEMENT'
	)
	)
	OR (
	owner = 'OPS$BRDB' and table_name in (
		'BRDB_ACC_NODE_PRODUCT_MAPPINGS', 
		'BRDB_APS_DELIVERY_FORMATS', 
		'BRDB_APS_MANUAL_TXN_AUDIT', 
		'BRDB_APS_MC_TXNS', 
		'BRDB_APS_MC_TXNS_E', 
		'BRDB_APS_PROD_SANITISE', 
		'BRDB_APS_RECON', 
		'BRDB_APS_REGION_MAPPING', 
		'BRDB_ARCHIVED_TABLES', 
		'BRDB_BRANCH_DECL', 
		'BRDB_BRANCH_DECL_ITEM', 
		'BRDB_BRANCH_FULL_EVENTS', 
		'BRDB_BRANCH_INFO', 
		'BRDB_BRANCH_NODE_INFO', 
		'BRDB_BRANCH_NODE_MONITOR', 
		'BRDB_BRANCH_ROLES', 
		'BRDB_BRANCH_ROLE_SERVICES', 
		'BRDB_BRANCH_STOCK_UNITS', 
		'BRDB_BRANCH_USERS', 
		'BRDB_BRANCH_USER_LAST_LOGON', 
		'BRDB_BRANCH_USER_POID_MAPPING', 
		'BRDB_BRANCH_USER_ROLES', 
		'BRDB_BRANCH_USER_SESSIONS', 
		'BRDB_CASH_DETAILS', 
		'BRDB_CASH_HEADER', 
		'BRDB_CLEARED_CLOSURE_DATA', 
		'BRDB_CLEARED_CONTROL_DATA', 
		'BRDB_COUNTER_MODE_CONVERSIONS', 
		'BRDB_CREDENCE_FILE_TOTALS', 
		'BRDB_CUTOFF_DETAILS', 
		'BRDB_CUTOFF_MARKERS', 
		'BRDB_CUTOFF_TOTALS', 
		'BRDB_DAILY_CUMULATIVE_SUMMARY', 
		'BRDB_DAILY_SUMMARY', 
		'BRDB_DESKTOP_MEMO_USER_DISTR', 
		'BRDB_DEVICE_TYPE_PRINCIPAL_MAP', 
		'BRDB_EXT_CASH_PRODUCTS', 
		'BRDB_EXT_COL_MAPPINGS', 
		'BRDB_EXT_ERROR_CODES', 
		'BRDB_EXT_FEED_REPORTS', 
		'BRDB_EXT_INTERFACE_FEEDS', 
		'BRDB_FAD_HASH_INSTANCE_MAPPING', 
		'BRDB_FAD_HASH_OUTLET_MAPPING', 
		'BRDB_FAD_HASH_VALUES', 
		'BRDB_FILES_TO_HOUSEKEEP', 
		'BRDB_F_HD_APS_TRANSACTIONS', 
		'BRDB_F_HD_DCS_TRANSACTIONS', 
		'BRDB_F_HD_EPOSS_EVENTS', 
		'BRDB_F_HD_EPOSS_TRANSACTIONS', 
		'BRDB_F_RX_APS_TRANSACTIONS', 
		'BRDB_F_RX_DCS_TRANSACTIONS', 
		'BRDB_F_RX_EPOSS_EVENTS', 
		'BRDB_F_RX_EPOSS_TRANSACTIONS', 
		'BRDB_F_RX_NWB_TRANSACTIONS', 
		'BRDB_F_ST_APS_TRANSACTIONS', 
		'BRDB_F_ST_DCS_TRANSACTIONS', 
		'BRDB_F_ST_EPOSS_EVENTS', 
		'BRDB_F_ST_EPOSS_TRANSACTIONS', 
		'BRDB_F_ST_MAL_TRANSACTIONS', 
		'BRDB_HOST_IF_FEEDS_MONITOR', 
		'BRDB_HOST_INTERFACE_FEEDS', 
		'BRDB_LAST_POST_INDICATOR', 
		'BRDB_NON_FI_RESP_CODES', 
		'BRDB_PASSWORD_HISTORY', 
		'BRDB_PED_KEYS_A', 
		'BRDB_PED_KEYS_B', 
		'BRDB_POID_CURRICULA', 
		'BRDB_POID_USER_DETAILS', 
		'BRDB_POUCH_COLL_DETAILS', 
		'BRDB_POUCH_COLL_HEADER', 
		'BRDB_POUCH_DEL_DETAILS', 
		'BRDB_POUCH_DEL_HEADER', 
		'BRDB_PROCESSES', 
		'BRDB_PS_BARCODES', 
		'BRDB_REM_OUT_POUCH_DETAILS', 
		'BRDB_REM_OUT_POUCH_HEADER', 
		'BRDB_RX_APS_TRANSACTIONS', 
		'BRDB_RX_BTS_DATA', 
		'BRDB_RX_BTS_DATA_ACC_CODE', 
		'BRDB_RX_BUREAU_TRANSACTIONS', 
		'BRDB_RX_CUT_OFF_SUMMARIES', 
		'BRDB_RX_DCS_TRANSACTIONS', 
		'BRDB_RX_EPOSS_EVENTS', 
		'BRDB_RX_EPOSS_TRANSACTIONS', 
		'BRDB_RX_GUARANTEED_REVERSALS', 
		'BRDB_RX_MESSAGE_JOURNAL', 
		'BRDB_RX_NRT_TRANSACTIONS', 
		'BRDB_RX_NWB_TRANSACTIONS', 
		'BRDB_RX_PBS_TRANSACTIONS', 
		'BRDB_RX_RECOVERY_TRANSACTIONS', 
		'BRDB_RX_REPORT_REPRINTS', 
		'BRDB_RX_REP_EVENT_DATA', 
		'BRDB_RX_REP_SESSION_DATA', 
		'BRDB_RX_TT_TRANSACTIONS', 
		'BRDB_RX_UNDO_TRANSACTIONS', 
		'BRDB_SQL_HINTS', 
		'BRDB_STOCK_UNIT_ASSOCIATIONS', 
		'BRDB_SUBPARTITION_RANGES', 
		'BRDB_SUB_FILE_AUDIT', 
		'BRDB_SUSPENDED_CUSTOMER_SESS', 
		'BRDB_SU_OPENING_BALANCE', 
		'BRDB_SU_PENDING_TRANSFER', 
		'BRDB_SU_PENDING_TRANSFER_DET', 
		'BRDB_SYSTEM_PARAMETERS', 
		'BRDB_TABLE_GROUPS', 
		'BRDB_TABLE_PARTITIONS', 
		'BRDB_TA_RECON', 
		'BRDB_TC_MODES_MAPPING', 
		'BRDB_TC_RECEIVED', 
		'BRDB_TXN_CORR_TOOL_CTL', 
		'BRDB_TXN_CORR_TOOL_JOURNAL', 
		'BRDB_UNAVAILABLE_SERVICES', 
		'C_BRDB_BRANCH_DECL', 
		'C_BRDB_BRANCH_DECL_ITEM', 
		'C_BRDB_BRANCH_FULL_EVENTS', 
		'C_BRDB_BRANCH_STOCK_UNITS', 
		'C_BRDB_BRANCH_USERS', 
		'C_BRDB_BRANCH_USER_LAST_LOGON', 
		'C_BRDB_BRANCH_USER_POID_MAPPI', 
		'C_BRDB_BRANCH_USER_ROLES', 
		'C_BRDB_BRANCH_USER_SESSIONS', 
		'C_BRDB_CASH_DETAILS', 
		'C_BRDB_CASH_HEADER', 
		'C_BRDB_CUTOFF_DETAILS', 
		'C_BRDB_CUTOFF_MARKERS', 
		'C_BRDB_CUTOFF_TOTALS', 
		'C_BRDB_DAILY_CUMULATIVE_SUMMAR', 
		'C_BRDB_DAILY_SUMMARY', 
		'C_BRDB_DESKTOP_MEMO_USER_DISTR', 
		'C_BRDB_LAST_POST_INDICATOR', 
		'C_BRDB_PASSWORD_HISTORY', 
		'C_BRDB_POUCH_COLL_DETAILS', 
		'C_BRDB_POUCH_COLL_HEADER', 
		'C_BRDB_POUCH_DEL_DETAILS', 
		'C_BRDB_POUCH_DEL_HEADER', 
		'C_BRDB_PS_BARCODES', 
		'C_BRDB_REM_OUT_POUCH_DETAILS', 
		'C_BRDB_REM_OUT_POUCH_HEADER', 
		'C_BRDB_RX_APS_TRANSACTIONS', 
		'C_BRDB_RX_BTS_DATA', 
		'C_BRDB_RX_BUREAU_TRANSACTIONS', 
		'C_BRDB_RX_CUT_OFF_SUMMARIES', 
		'C_BRDB_RX_DCS_TRANSACTIONS', 
		'C_BRDB_RX_EPOSS_EVENTS', 
		'C_BRDB_RX_EPOSS_TRANSACTIONS', 
		'C_BRDB_RX_GUARANTEED_REVERSALS', 
		'C_BRDB_RX_MESSAGE_JOURNAL', 
		'C_BRDB_RX_NRT_TRANSACTIONS', 
		'C_BRDB_RX_NWB_TRANSACTIONS', 
		'C_BRDB_RX_PBS_TRANSACTIONS', 
		'C_BRDB_RX_RECOVERY_TRANSACTION', 
		'C_BRDB_RX_REPORT_REPRINTS', 
		'C_BRDB_RX_REP_EVENT_DATA', 
		'C_BRDB_RX_REP_SESSION_DATA', 
		'C_BRDB_RX_TT_TRANSACTIONS', 
		'C_BRDB_RX_UNDO_TRANSACTIONS', 
		'C_BRDB_STOCK_UNIT_ASSOCIATIONS', 
		'C_BRDB_SUSPENDED_CUSTOMER_SESS', 
		'C_BRDB_SU_OPENING_BALANCE', 
		'C_BRDB_SU_PENDING_TRANSFER', 
		'C_BRDB_SU_PENDING_TRANSFER_DET', 
		'C_BRDB_TXN_CORR_TOOL_JOURNAL', 
		'C_LFS_PLO_DETAILS', 
		'C_LFS_PLO_HEADER', 
		'C_LFS_RDC_DETAILS', 
		'C_LFS_RDC_HEADER', 
		'C_RDDS_DESKTOP_MEMO_DISTR', 
		'C_TPS_TXN_ACK_DETAILS', 
		'C_TPS_TXN_CORRECTION_DETAILS', 
		'LFS_PLO_DETAILS', 
		'LFS_PLO_HEADER', 
		'LFS_RDC_DETAILS', 
		'LFS_RDC_HEADER', 
		'PATCH_RELEASE_MANAGEMENT', 
		'RDDS_ACCOUNTING_NODES', 
		'RDDS_APS_CLIENT_ACCOUNTS', 
		'RDDS_APS_DELIVERY_AGREEMENTS', 
		'RDDS_APS_REGION_BANK_HOLIDAYS', 
		'RDDS_AP_TOKENS', 
		'RDDS_BANK_HOLIDAYS', 
		'RDDS_BRANCHES', 
		'RDDS_BRANCH_BUREAU_REGION', 
		'RDDS_BRANCH_OPENING_PERIODS', 
		'RDDS_CHECKSUM', 
		'RDDS_CHECKSUM_HIST', 
		'RDDS_CLIENTS', 
		'RDDS_CLIENT_ACCOUNTS', 
		'RDDS_DELIVERY', 
		'RDDS_DELIVERY_TYPE', 
		'RDDS_DESKTOP_MEMO', 
		'RDDS_DESKTOP_MEMO_DISTR', 
		'RDDS_DEVICE_DELIVERY', 
		'RDDS_DEVICE_PACKAGE_CONTENT', 
		'RDDS_LOGON_CURRICULA_DETAILS', 
		'RDDS_LOGON_CURRICULA_GROUPS', 
		'RDDS_PACKAGE', 
		'RDDS_PACKAGE_CONTENT', 
		'RDDS_PACKAGE_TYPE', 
		'RDDS_POUCH_TYPES', 
		'RDDS_PRODUCTS', 
		'RDDS_PRODUCT_CURRENCY', 
		'RDDS_PRODUCT_GROUPS', 
		'RDDS_PRODUCT_MODES', 
		'RDDS_PS_PRODUCT_MAP', 
		'RDDS_RETAIL_PACKAGE', 
		'RDDS_TANDT_SERVICE_RULES', 
		'RDDS_TRADE_CURRICULA_DETAILS', 
		'RDDS_TRADE_CURRICULA_GROUPS', 
		'RDDS_TRANSMISSION_SOURCE', 
		'RDDS_TRANS_MODES', 
		'RD_PACKAGE_DELIVERY', 
		'REPL_TRANS_LATEST', 
		'TPS_TXN_ACK_DETAILS', 
		'TPS_TXN_CORRECTION_DETAILS'
	)
	
	)
	order  by owner, table_name, column_id
	) WHERE lower(table_name) NOT IN ('import_csv_bulk_file_audit','obc_stock_unit_type','obc_interface_audit','emdb2_associated_stock_unit','c_lfs_rdc_header','c_brdb_desktop_memo_user_distr','c_brdb_device_type_principal_map','c_brdb_ext_feed_reports','c_brdb_ped_keys_a','c_brdb_ped_keys','c_brdb_ped_keys_b','c_brdb_rx_message_journal','c_brdb_rx_report_reprints','c_brdb_branch_decl','brdb_desktop_memo_user_distr','brdb_device_type_principal_map','brdb_ext_feed_reports','brdb_ped_keys_a','brdb_ped_keys','brdb_ped_keys_b','brdb_rx_message_journal','brdb_rx_report_reprints')
	
	) GROUP BY table_name,owner